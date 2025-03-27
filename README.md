# medflow [![Test](https://github.com/dyammarcano/medflow/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/dyammarcano/medflow/actions/workflows/test.yml)

PoC de un sistema de flujo de pacientes en un hospital, con servicios de registro, evaluación médica, internación y
monitoreo central.

## 🏥 Flujo de Entrada del Paciente:

|           Etapa           |       Nombre del Servicio        |        Topic NATS        |
|:-------------------------:|:--------------------------------:|:------------------------:|
|      Entrada inicial      |     	patient-intake-service	     | operation.incoming.data  |
|         Registro	         |  patient-registration-service	   |  operation.stage1.data   |
|    Evaluación médica	     |   medical-assessment-service	    |  operation.stage2.data   |
| Internación / Aprobación	 | admission-authorization-service	 |  operation.stage3.data   |
|    Casos prioritarios	    |  priority-case-handler-service	  | operation.priority1.data |

## 🧪 Servicios de Examen:

|   Tipo de Examen   | Nombre del Servicio |      Topic NATS       |
|:------------------:|:-------------------:|:---------------------:|
| Exámenes generales | general-lab-service | operation.examas.data |
|    Imagenología    | imaging-lab-service | operation.examas.data |
|   Especializados   | special-lab-service | operation.examas.data |

## 📡 Monitor Central:

|       Función       |   Nombre del Servicio    |    Topic NATS    |
|:-------------------:|:------------------------:|:----------------:|
| Monitor y dashboard | clinical-monitor-service | operation.*.data |

## Temas (Subjects) definidos en NATS:

operation.incoming.data: Entrada inicial del paciente.

operation.stage1.data hasta operation.stage3.data: Servicios de procesamiento por etapas.

operation.examas.data: Eventos para exámenes (todos los servicios de examen escuchan este subject).

operation.response.data: Todos los servicios de examen envían su resultado aquí.

operation.error.data: Mensajes con errores en cualquier etapa.

## Roles de cada componente:

Servicios de etapa (stage1, stage2, etc):
Escuchan su subject (por ejemplo: operation.stage1.data).

Procesan el evento.

Publican al siguiente subject (ej. operation.stage2.data).

Persisten localmente en SQLite.

Servicios de exámenes:
Todos escuchan operation.examas.data.

Procesan sólo si el evento coincide con su exam_type.

Publican el resultado al monitor vía operation.response.data.

Monitor:
Se suscribe a operation.*.data.

Guarda todos los eventos en PostgreSQL.

Se suscribe a operation.response.data para guardar respuestas de exámenes.

Expone WebSocket para frontend con checklist en tiempo real.

## 🔧 Paso 1: Ejecutar PostgreSQL con Podman

```bash
podman run --name medflow-postgres -d \
-e POSTGRES_USER=user \
-e POSTGRES_PASSWORD=password \
-e POSTGRES_DB=medflow \
-p 5432:5432 \
docker.io/postgres:latest
```

## ✅ Paso 2: Ejecutar NATS JetStream con Podman

```bash
podman run --name nats-jetstream -d \
-p 4222:4222 \
-p 8222:8222 \
docker.io/nats:latest \
-js -D -m 8222
```