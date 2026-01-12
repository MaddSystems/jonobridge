Componente 1: Chequeo de Geocercas (Expresiones Nativas)

   * Nombre: CheckGeofenceExclusion
   * Qué hace: Verifica que la ubicación actual del vehículo NO esté dentro de ninguna de las tres zonas de exclusión ("Taller", "CLIENTES",
     "Resguardo Cedis Puerto").
   * Cómo lo hace: Utiliza las funciones de geocercas nativas de GPSGate.
   * Equivalente en Grule: Ya lo tenemos con !state.IsInsideGroup(...).

  ---

  Componente 2: Chequeo de Conectividad (Expresión Nativa)

   * Nombre: CheckOfflineStatus
   * Qué hace: Confirma que el vehículo ha estado "offline" (sin reportar una posición válida) durante al menos 5 minutos.
   * Cómo lo hace: Utiliza una función de tiempo nativa de GPSGate.
   * Equivalente en Grule: Ya lo tenemos con state.IsOfflineFor(5).

  ---

  Componente 3: Gestión del Buffer de Memoria (Script)

   * Nombre: UpdateCircularBuffer
   * Qué hace: SIEMPRE, en cada ejecución del script, añade la trama actual a un buffer en memoria y se asegura de que el buffer nunca exceda las
     10 posiciones (elimina la más antigua si es necesario).
   * Cómo lo hace:
       1. Obtiene el estado context.state (nuestro state en Grule).
       2. Si no existe, lo inicializa con la trama actual.
       3. Si existe, añade la nueva trama (.push(actualData)).
       4. Si el tamaño supera 10, elimina la más antigua (.shift()).
       5. Establece una bandera doit = true solo cuando el buffer está lleno (contiene 10 posiciones).
   * Equivalente en Grule: Ya lo tenemos con state.UpdateMemoryBuffer(...).

  ---

  Componente 4: Condición Principal de Evaluación (Script)

   * Nombre: EvaluateOnlyOnInvalid
   * Qué hace: Es el "guardián" de la lógica de detección. Solo permite que el código de análisis de métricas se ejecute si se cumplen dos
     condiciones:
       1. La trama actual es inválida (!isValid).
       2. El buffer de memoria está lleno (doit es true).
   * Cómo lo hace: Con la estructura if(!isValid){ if (doit) { ... } }.
   * Equivalente en Grule: Esto es lo que debemos replicar en una única condición when.

  ---

  Componente 5: Cálculo de Métricas del Buffer (Script)

   * Nombre: CalculateBufferMetrics
   * Qué hace: Si el Componente 4 lo permite, itera sobre las 10 tramas del buffer para calcular dos métricas clave:
       1. Velocidad Promedio: Suma la velocidad de todas las tramas dentro de una ventana de tiempo (maximumtime = 90 minutos) y las divide por el
          número de tramas consideradas.
       2. Señal GSM Promedio: Suma el nivel de señal GSM de las últimas 5 tramas y lo divide entre 5.
   * Cómo lo hace: Un bucle .forEach que recorre info['actualData'] y acumula los valores de km y signalLevel en variables.
   * Equivalente en Grule: Ya lo tenemos con state.CalculateJammerMetricsIfReady(), pero debemos asegurarnos de llamarlo solo cuando se cumplan las
     condiciones del Componente 4.

  ---

  Componente 6: Evaluación de Umbrales (Script)

   * Nombre: CheckMetricThresholds
   * Qué hace: Compara las métricas calculadas en el Componente 5 contra los umbrales definidos:
       1. Velocidad promedio (kilometrajePromedio) debe ser >= 10.
       2. Señal GSM promedio (senalPromedio) debe ser >= 9.
   * Cómo lo hace: Un if(kilometrajePromedio>=10){ if(senalPromedio>=9){ return true; } }.
   * Cómo se integra: El return true de este bloque es lo que hace que la expresión del "Script" completa sea verdadera, lo que, junto con las
     expresiones nativas, dispara la alerta.
   * Equivalente en Grule: Ya lo tenemos con state.GetCounter("jammer_avg_speed_90min") >= 10 y state.GetCounter("jammer_avg_gsm_last5") >= 9.

  ---

  Basado en tu razonamiento, aquí está la estructura lógica que integra todos los componentes en el orden más eficiente, replicando el flujo de
  GPSGate:

   1. Componente 3 (`UpdateCircularBuffer`): Siempre se ejecuta primero, sin condiciones. Es una acción, no una evaluación. Mantiene el buffer de
      10 posiciones actualizado en todo momento. Esto lo haremos en una regla separada con alta prioridad (salience), asegurando que ocurra antes
      que cualquier evaluación.

   2. Componente 4 (`EvaluateOnlyOnInvalid`): Esta es la puerta de entrada principal para la detección. Es el primer chequeo en nuestra regla
      principal. Si la trama actual es VÁLIDA, la regla entera se detiene aquí y no se evalúa nada más. Esto ahorra muchísimo procesamiento.

   3. Componente 2 (`CheckOfflineStatus`): Si la trama es inválida, lo siguiente más rápido de chequear es el estado "offline". Es un chequeo de
      tiempo en memoria, muy eficiente. Si el vehículo no ha estado offline por 5 minutos, no seguimos.

   4. Componente 5 (`CalculateBufferMetrics`): Solo si la trama es inválida y el vehículo está offline, procedemos a calcular las métricas del
      buffer. Es importante hacerlo aquí, porque las siguientes condiciones dependen de estos valores.

   5. Componente 6 (`CheckMetricThresholds`): Con las métricas recién calculadas, hacemos la comparación numérica de umbrales (velocidad >= 10, gsm
      >= 9). Es una operación matemática simple y rápida.

   6. Componente 1 (`CheckGeofenceExclusion`): Este es el último y más "costoso" chequeo. Solo si la trama es inválida, el vehículo está offline, Y
      las métricas del buffer cumplen los umbrales, nos tomamos la molestia de consultar la base de datos para verificar las geocercas.

  Cómo encaja todo en la(s) regla(s) de Grule:

  