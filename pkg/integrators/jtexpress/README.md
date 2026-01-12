# Jtexpress 

Proceso de Consumo e Integración del Servicio API

El consumo de los servicios de la plataforma Open API se realiza principalmente a través del proceso de pruebas de comunicación, validación de parámetros y generación de firmas digitales, lo cual ocurre después de que la cuenta ha sido registrada y aprobada.
I. Preparación y Autenticación
Para asegurar una integración correcta antes de realizar la conexión productiva, el cliente debe validar la autenticación y la estructura de los datos:

1. Generación de Firmas Digitales (Digest): El apartado "Open API platform – SDK" permite generar y validar las firmas digitales (digest) necesarias para utilizar los servicios de la API de J&T.
    ◦ En Business parameter signature se ingresan credenciales para obtener el digest de negocio.
    ◦ En Headers signature se generan las firmas de encabezado.
    ◦ El SDK sirve como herramienta de prueba para asegurar que los parámetros y autenticaciones estén correctos antes de la integración real.
2. Uso de Credenciales de Prueba: Para el proceso de pruebas, se deben utilizar las credenciales específicas de prueba proporcionadas en la documentación, como apiAccount de prueba: 292508153084379141 y privateKey de prueba: a0a1047cce70493c9d5d29704f05d0d9.
II. Pruebas de Comunicación
Una vez que la plataforma está certificada y se tienen las firmas listas, se procede a la prueba de los servicios:
1. Ubicación de los Servicios: Los servicios que el cliente desea aplicar se visualizan en el menú de "Documentación de la API".
2. Ejecución de Pruebas: Las pruebas de comunicación se pueden realizar en el apartado de "Consola" (https://open.jtjms-mx.com/#/control), donde se encuentran las URLs de los servicios de J&T utilizables.
3. Requisito Mínimo: Para cada servicio que se desee utilizar, es necesario realizar al menos tres (3) pruebas de comunicación exitosas.
4. Validación de Parámetros:
    ◦ Para cada servicio, se detallan los parámetros que se deben incluir para una comunicación exitosa en la parte inferior de cada web service.
    ◦ Los parámetros que contienen una 'Y' son requeridos; los que contienen una 'N' no son necesarios y pueden ser omitidos.
    ◦ El apartado "TEST TOOLS" incluye información de ejemplo de los datos que debe contener cada servicio, además de los códigos de error.
III. Activación y Uso Productivo
1. Solicitud de Aplicación ("Application"): Si las tres pruebas por servicio son exitosas, al dar click en "Application", el servicio se manda a revisión y aparecerá como "Applying".
2. Aprobación del Servicio: Inicialmente, los servicios aparecen como "Not coordinated". Cuando los servicios se aprueben, cambiarán su estado a "Success".
3. Gestión de Interfaz: En el apartado número 5, "Interface management", se puede solicitar la revisión del servicio para poder hacer uso de las URLs de forma productiva.
4. Servicios de Retorno Específicos: Para los servicios Logistics trajectory return y Order status return, es obligatorio ingresar las URLs a las que se realizará el retorno de la información o la respuesta del servicio.
