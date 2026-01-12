Pasar el objeto como un puntero, resolvio el problema de la "desconexión de memoria" que causaba el bucle infinito.

1. El Fix de Memoria

Go

// Antes: dataContext.Add("Jono", wrapper)  <- Copiaba los datos
// Ahora:
err = dataContext.Add("Jono", &wrapper) // <- Pasa la dirección real
Por qué esto detiene el loop: Cuando la regla en el archivo .grl ejecuta Jono.Cond3Passed = true, ahora está escribiendo directamente en la dirección de memoria de wrapper.

Sin el &, Grule modificaba una copia interna; entonces, cuando se llamaba a Changed(), el motor volvía a leer tu variable original (que seguía siendo false) y se disparaba otra vez. Con el puntero, en la segunda vuelta del motor, este ve que Cond3Passed ya es true, por lo que el when de la regla ya no se cumple y el ciclo se rompe naturalmente.

2. Eliminación de Changed() vs Punteros
Como ahora usas un puntero, tienes dos opciones:

Si dejas el Changed(): El motor reiniciará la evaluación, pero solo una vez, porque encontrará que la condición ya se marcó como true.

Si quitas el Changed(): El motor seguirá evaluando las reglas hacia abajo (por salience) sin reiniciar todo el árbol. Esto es más eficiente (consume menos CPU).

3. Una observación en detectExecutedComponents
He notado que en la función detectExecutedComponents, ahora recibes wrapper por valor:

Go

func detectExecutedComponents(wrapper PacketWrapper) []string { ... }
Esto está bien para esa función, ya que solo necesitas leer los datos para el log de auditoría una vez que el motor terminó. No afecta al loop.

4. ¿Qué esperar en los logs ahora?
Si ejecutas el comando grep que usaste antes: grep -c "successfully selected rule candidate for execution after 100 cycles" logs.txt

Debería devolver 0. Si el número sigue subiendo, significa que hay otra regla (quizás en DEFCON 1 o 2) donde te falta la validación !CondXPassed en el when.



Esta nueva versión de la simulación y la configuración del motor están perfectamente alineadas. Has resuelto el problema de la "historia" del dispositivo que impedía que la lógica avanzara.

Aquí te detallo por qué esta implementación sí debería completar "la película" hasta el envío del Telegram:

1. El "Calentamiento" del Buffer (Fase 1)
El cambio en send_periodic4defcon4.py es la clave del éxito. Al enviar 11 paquetes con status="A", garantizas que:


BufferHas10 sea ✅: El motor de reglas ahora verá un historial confiable antes de que empiece el ataque.



Métricas Iniciales: Se establece un promedio de velocidad de 50 km/h y señal GSM de 20, lo que deja el sistema listo para cumplir los umbrales de la Fase 2.


2. Superación del Umbral Temporal (Fase 2)
Al enviar 24 paquetes inválidos (status="V") con intervalos de 15 segundos (6 minutos en total), aseguras que la condición Jono.IsOfflineFor5Min se cumpla sin dudas.


Como el buffer ya estaba lleno desde la Fase 1, el motor no se detendrá por falta de datos y escalará de DEFCON 1 a DEFCON 2 en cuanto el reloj lógico supere los 300 segundos.

3. Integridad en el grule_worker.go
La implementación del PacketWrapper en el worker ahora captura correctamente el estado de CurrentlyInvalid. Esto es vital porque:

El worker identifica que la trama es válida para procesamiento pero inválida en posición (Status V), activando la lógica de "pérdida de contacto" en el motor de reglas.

La copia por valor (localPacket := packet) evita que el estado de un paquete de calentamiento se mezcle con uno de ataque.

4. Flujo de Alerta (DEFCON 4)
Una vez que se cumplen los 5 minutos y los promedios se mantienen (Speed > 10 y GSM > 9), la regla DEFCON4_JammerDetected finalmente se activará:

Se enviará el mensaje a Telegram con los datos reales del ataque.

El flag AlertFired se pondrá en true, cerrando el ciclo de detección.

Siguiente paso recomendado: Ejecuta de nuevo tu script de análisis analyze_jammer_flow.py. Ahora deberías ver cómo la columna B10 se marca con ✅ desde las primeras filas y cómo, tras llegar al minuto 5 de la Fase 2, las columnas MET y STAGE avanzan rápidamente hasta defcon4_alert.

