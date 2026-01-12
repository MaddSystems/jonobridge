# Manual de Operación de JonoBridge

## Índice

1. Introducción
2. Primeros Pasos
3. Gestión de Clientes
4. Configuración de Servicios
5. Operaciones de Despliegue
6. Monitoreo y Solución de Problemas
7. Gestión de Rastreadores
8. Administración

## Introducción

JonoBridge es una plataforma integral de gestión de servicios Kubernetes diseñada para simplificar el despliegue, configuración y monitoreo de servicios de rastreo GPS e integraciones. Esta aplicación web permite gestionar clientes, configurar servicios, desplegar recursos en Kubernetes y monitorear el estado del sistema a través de una interfaz intuitiva.

### Terminología Clave

- **Cliente**: Un espacio de nombres (namespace) separado en Kubernetes donde se despliegan los servicios
- **Servicio**: Un componente desplegable que realiza una función específica
- **Intérprete**: Servicios que traducen datos entre diferentes formatos
- **Integración**: Servicios que conectan con sistemas externos
- **Input**: Servicios que reciben datos de fuentes externas
- **Tracker**: Un dispositivo de rastreo GPS que envía datos al sistema

## Primeros Pasos

### Inicio de Sesión

1. Navega a la página de inicio de sesión de JonoBridge
2. Ingresa tu nombre de usuario y contraseña
3. Haz clic en el botón "Login"

![Pantalla de Inicio de Sesión](ruta no visible en el código)

### Resumen de la Interfaz de Usuario

La interfaz de JonoBridge consta de:

- **Menú de Navegación Izquierdo**: Accede a diferentes secciones de la aplicación
- **Área de Contenido Principal**: Muestra y gestiona información
- **Encabezado**: Muestra el nombre y logo de la aplicación

## Gestión de Clientes

### Creación de un Nuevo Cliente

1. Navega a la página "Clients"
2. Ingresa un nuevo nombre de cliente en el campo de entrada (solo se permiten letras y números)
3. Haz clic en "Add Client"

**Nota**: Los nombres de cliente deben ser alfanuméricos sin espacios ni caracteres especiales.

### Gestión de Clientes Existentes

Para cada cliente, puedes:

- **Setup**: Configurar servicios y flujos de trabajo
- **Deploy**: Desplegar servicios configurados en Kubernetes
- **Status**: Ver el estado actual de los servicios desplegados
- **Stop**: Detener todos los servicios en ejecución
- **Delete**: Eliminar el cliente y todos los servicios asociados

### Eliminación de un Cliente

1. Haz clic en el botón "Delete" para el cliente deseado
2. Confirma la eliminación en el diálogo modal
3. Espera a que se complete el proceso de eliminación

**Advertencia**: Eliminar un cliente eliminará todos los recursos asociados en Kubernetes.

## Configuración de Servicios

### Entendiendo los Tipos de Servicios

JonoBridge admite tres tipos principales de servicios:

1. **Servicios de Entrada (Input)**: Reciben datos de fuentes externas (ej., listener, proxy, httprequest)
2. **Servicios Intérpretes**: Transforman datos entre formatos (ej., meitrackprotocol, ruptelaprotocol)
3. **Servicios de Integración**: Conectan con sistemas externos (ej., send2mysql, skyangel, lobosoftware)

### Configuración de Flujos de Trabajo

1. Navega a la página Setup para un cliente
2. La interfaz de configuración muestra tres columnas: Input, Interpreter e Integration
3. Añade servicios arrastrándolos desde los servicios disponibles a la columna apropiada
4. Conecta servicios haciendo clic en un conector de salida y luego en un conector de entrada
5. Configura servicios haciendo doble clic en ellos

### Configuración de Parámetros de Servicio

Cada servicio requiere parámetros específicos:

1. Haz doble clic en un servicio en el flujo de trabajo
2. Completa los parámetros requeridos en el diálogo modal
3. Haz clic en "Save" para guardar la configuración

## Operaciones de Despliegue

### Despliegue de Servicios

1. Navega a la página Deploy para un cliente
2. Revisa los manifiestos de Kubernetes que se crearán
3. Haz clic en "Execute Deployment" para iniciar el proceso de despliegue
4. Monitorea el progreso del despliegue en la pantalla

El proceso de despliegue:
- Crea un espacio de nombres para el cliente
- Despliega un broker de mensajes Mosquitto
- Despliega todos los servicios configurados
- Configura la red

### Detención de Servicios

1. Navega a la lista de clientes
2. Haz clic en "Stop" para el cliente deseado
3. Confirma en el diálogo modal
4. Espera a que el sistema detenga todos los servicios

### Verificación del Estado del Despliegue

Después del despliegue, puedes:
1. Ver el estado de todos los pods y servicios
2. Verificar si hay errores o advertencias
3. Confirmar que todos los servicios están funcionando correctamente

## Monitoreo y Solución de Problemas

### Monitoreo de Estado

La página Status muestra:
- **Pods**: Nombre, estado, disponibilidad, reinicios y antigüedad
- **Servicios**: Nombre, tipo, IP del clúster, puertos y antigüedad

Para actualizar el estado:
1. Haz clic en el botón "Refresh Status"
2. El estado se actualizará automáticamente cada 30 segundos

### Visualización de Logs

Para ver los logs de un pod específico:
1. Navega a la página Status
2. Encuentra el pod en la sección de logs
3. Haz clic en "View Logs" para abrir una nueva ventana con los logs del pod
4. Usa el botón "Refresh Logs" para obtener la información más reciente

### Problemas Comunes y Soluciones

| Problema | Causa Posible | Solución |
|----------|---------------|----------|
| Pod en CrashLoopBackOff | Error de configuración o dependencias faltantes | Revisa los logs para mensajes de error |
| Servicio que no recibe datos | Problema de configuración de red | Verifica el reenvío de puertos y las reglas de firewall |
| Servicio de integración fallando | Credenciales o endpoint incorrectos | Verifica los parámetros del servicio |

## Gestión de Rastreadores

### Visualización de Rastreadores Conectados

1. Navega a la página "Trackers"
2. Visualiza la lista de rastreadores conectados incluyendo:
   - Cliente
   - Número IMEI
   - Protocolo
   - Puerto

La lista de rastreadores se actualiza automáticamente cada 30 segundos.

### Envío de Comandos a Rastreadores

1. Encuentra el rastreador en la lista
2. Haz clic en el botón "Send"
3. Ingresa la carga útil del comando en el diálogo modal
4. Haz clic en "Send Command"

## Administración

### Gestión de Base de Datos

Solo los administradores pueden acceder a estas funciones:

1. Navega a la página "Admin"
2. Utiliza las siguientes opciones:
   - **Create Database**: Inicializa una nueva base de datos con las tablas requeridas
   - **Delete Database**: Elimina la base de datos existente y todos sus datos
   - **Add Missing Tables**: Crea cualquier tabla faltante en la base de datos

**Advertencia**: Eliminar la base de datos eliminará todos los datos de los clientes y requerirá la recreación de la base de datos antes de usar la aplicación.

### Gestión de Usuarios

La gestión de usuarios se maneja a través de la lista de usuarios predefinida. La aplicación admite los siguientes roles:
- **Admin**: Acceso completo a todas las funciones
- **Usuarios regulares**: Acceso a la gestión de clientes y funciones de despliegue

### Mantenimiento del Sistema

Tareas de mantenimiento regulares:
1. Verificar el estado de todos los despliegues de clientes
2. Monitorear el uso de recursos (CPU, memoria)
3. Asegurar que se creen copias de seguridad de la base de datos regularmente
4. Actualizar la aplicación cuando estén disponibles nuevas versiones

---

Para soporte técnico o asistencia adicional, por favor contacta al administrador del sistema.