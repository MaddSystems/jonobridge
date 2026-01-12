**SERVICIOS REST SAI-TI GPS - \(DICSA\)** **End Points **

https://viatix.com.mx/red\_energy/API\_GPS/data/\{\{params\}\}/ 

https://viatix.com.mx/red_energy/AP/_GPS/data/{{params}}/

**GetToken **

**GET: ** {{URL}}/token 

https://viatix.com.mx/red_energy/API_GPS/token/?user=dicuser&password=dicpwd25



**Credenciales de acceso: **

Usuario: dicuser 

Password: dicpwd25 



**Instrucciones de consumo. **

**Nota: **

**Puede utilizar POSTMAN para realizar las pruebas de consumo. **

**1.- Generar token de autenticación. **

En la siguiente pantalla se muestra como y donde se debe colocar la URL para generar el Token de acceso y la respuesta que se debe recibir. 

Tipo de método: GET 



**Parámetros \(Obligatorios y opcionales\)** 

User: Obligatorio 

Password: Obligatorio 

https://viatix.com.mx/red_energy/API_GPS/token/?user=dicuser&password=dicpwd25


Return Info:
{"token_access":"04fd08e7f6e61fddf1f1"}



**2.- Envío de parámetros. **

**Parámetros \(Obligatorios y opcionales\)** 

User: Obligatorio 

Password: Obligatorio 

Token: Obligatorio 

Placa: Obligatorio 

Neconomico: Opcional 

Latitud: Obligatorio 

Longitud: Obligatorio 

Altitud: Opcional 

Dirección: Opcional 

Velocidad: Opcional 

fecha_localizacion: Obligatorio 

edo_unidad: Opcional 






En la siguiente pantalla se muestra un ejemplo del envío de parámetros por el método POST. 

https://viatix.com.mx/red_energy/API_GPS/data/?user=dicuser&password=dicpwd25&token=3cb6293fe89ad68d26f78placa&neconomico&latitud&.....


| **Key**              | **Value**                |
|----------------------|--------------------------|
| user                 | dicuser                  |
| password             | dicpwd25                 |
| token                | 04fd08e7f6e61fddf1f1     |
| placa                |                          |
| neconomico           |                          |
| latitud              |                          |
| longitud             |                          |
| altitud              |                          |
| direccion            |                          |
| velocidad            |                          |
| fecha_localizacion   |                          |

Nota importante: 

El token tiene que ser enviado como se muestra en la pantalla anterior, es decir, se tiene que enviar como parámetros y también en el HEADER como se muestra a continuación: 

Headers:

| **Key**              | **Value**                |
|----------------------|--------------------------|
| Authorization        | 04fd08e7f6e61fddf1f1     |


Si todo está correcto entonces recibirá una respuesta de registro correcto: 



**Manejo de errores: **





Los errores que va a mandar la aplicación dependerá meramente de los campos que son obligatorios, así mismo si el token no ha sido enviado tanto en Header como en parámetros y/o el token ha caducado. 



Importante: El tiempo de vida del Token está parametrizado para que dure vigente 24 horas.



