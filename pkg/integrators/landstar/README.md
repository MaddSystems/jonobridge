# Landstar 

Dado que el servicio de geolocalización de Landstar Metro es de tipo SOAP, las solicitudes se realizan enviando mensajes XML específicos (envolventes SOAP) a través de HTTP POST. La prueba debe constar de dos pasos esenciales:
1. Autenticación: Obtener el token de sesión usando el método GetUserToken.
2. Geolocalización: Enviar los datos de ubicación usando el método GPSAssetTracking, utilizando el token obtenido.
Este código utiliza la librería estándar de Python requests para interactuar con el servicio.


Notas Importantes sobre la Prueba
1. Credenciales Obligatorias: El script está listo para ser usado, pero no funcionará hasta que reemplace los placeholders (SU_USUARIO_AQUI, SU_CONTRASENA_AQUI) con las credenciales reales que le debe proporcionar su contacto en Landstar Metro para el ambiente de pruebas.
2. Ambiente: El código utiliza la URL del Ambiente de Pruebas: https://compassvmqa.centralus.cloudapp.azure.com/locations/locationReceiver.wsdl.
3. Flujo Lógico: La prueba sigue el flujo definido en la documentación: obtener el token con usuario y contraseña (GetUserToken), y luego utilizar ese token en la petición de geolocalización (GPSAssetTracking).
4. Uso de XML: Las estructuras XML de las peticiones (soap_request) están copiadas directamente de los ejemplos proporcionados en los fuentes, asegurando la compatibilidad con el formato SOAP requerido.