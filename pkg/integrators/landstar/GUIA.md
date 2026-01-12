# GUIA

La guía que usted proporciona es un Contrato de Integración asociado al consumo del Servicio Web de Geolocalización de Unidades de Asociados Landstar Metro. Este documento es de naturaleza técnica y describe las especificaciones necesarias para la integración.
A continuación, se presenta un resumen de los puntos clave de la guía:
1. Propósito y Modalidad del Servicio
• Objetivo: El documento detalla el contrato de integración para consumir el Servicio Web de geolocalización de una unidad móvil (tracto camión) por medio de un GPS.
• Modalidad: El servicio opera bajo la modalidad SOAP (Simple Object Access Protocol).
• Alcance: El enfoque es técnico, detallando los aspectos relacionados con el procedimiento de integración, los mecanismos de consumo, los formatos a utilizar, los parámetros esperados y las posibles respuestas o excepciones de error.
2. Ambientes de Trabajo y Credenciales
Se definen dos ambientes de trabajo para la interacción:
Ambiente
URL (WSDL)
Credenciales
Pruebas (QA)
https://compassvmqa.centralus.cloudapp.azure.com/locations/locationReceiver.wsdl
Proporcionadas por su contacto en Landstar Metro.
Productivo
https://compass-landstar.centralus.cloudapp.azure.com/locations/locationReceiver.wsdl
Se proporcionan después de una prueba exitosa; se da un usuario y contraseña para cada flota de tracto camiones.
3. Métodos del Servicio Web
La integración se basa en dos métodos principales que deben consumirse secuencialmente:
A. GetUserToken (Autenticación)
• Función: Es el método utilizado para la autentificación.
• Requisito: Se debe proporcionar el Usuario (userId) y la Contraseña (password) suministrados por Landstar.
• Respuesta Exitosa: Devuelve un Token de sesión (token). Este token es crucial, ya que es solicitado para el consumo del siguiente servicio (GPSAssetTracking).
• Ejemplo de Error: Se muestra un error con el código 2004 y el mensaje "Usuario o contraseña son incorrectos".
B. GPSAssetTracking (Envío de Geolocalización)
• Función: Procesa los eventos de geolocalización recibidos desde el prestador una vez que se han obtenido las credenciales y el token.
• Requisito Clave: Se solicita el token previamente obtenido.
• Datos Obligatorios del Evento (<loc:Event>):
    ◦ Código (Event.Code): Equivalente al código del evento del AVL (ej. 0 para Sin Evento, 911 para Botón de pánico).
    ◦ Fecha / Hora (Event.Date): Fecha del evento en formato UTC (YYYY-MM-DTHH:MM:SS).
    ◦ Latitud (Event.latitude) y Longitud (implícito en la tabla, y requerido por ejemplos).
    ◦ Placa (Event.Asset): Sin espacios ni caracteres especiales.
    ◦ Dirección (Event.Direction): Dirección actual en texto.
    ◦ Velocidad (Event.Speed).
• Datos Opcionales: Incluyen la altitud, número de serie del GPS, identificadores del transportista y del viaje (Event.Shipment), odómetro, estado de ignición (True/False), nivel de batería, curso (Norte, Sur, Este, etc.) y valores de temperatura/humedad.
• Respuesta Exitosa: Devuelve un ID de trabajo (idJob).
• Ejemplos de Errores: Se especifican respuestas de error cuando faltan elementos obligatorios (como el code o la latitude).