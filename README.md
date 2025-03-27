# medflow [![Test](https://github.com/dyammarcano/medflow/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/dyammarcano/medflow/actions/workflows/test.yml)

## 🏥 Flujo de Entrada del Paciente:

|           Etapa            |        Nombre del Servicio         |         Topic NATS         |
|:----------------------------:|:----------------------------------:|:--------------------------:|
|      Entrada inicial       |  	patient-intake-service	   |  operation.incoming.data   |
|         Registro	          |   patient-registration-service	    |   operation.stage1.data    |
|     Evaluación médica	     |    medical-assessment-service	     |   operation.stage2.data    |
| Internación / Aprobación	  |  admission-authorization-service	  |   operation.stage3.data    |
|    Casos prioritarios	     |   priority-case-handler-service	   |  operation.priority1.data  |

## 🧪 Servicios de Examen:

| Tipo de Examen | Nombre del Servicio | Topic NATS |
|:----------------:|:---------------------:|:------------:|
| Exámenes generales | general-lab-service | operation.examas.data |
| Imagenología | imaging-lab-service | operation.examas.data |
| Especializados | special-lab-service | operation.examas.data |

## 📡 Monitor Central:

| Función | Nombre del Servicio | Topic NATS |
|:---------:|:---------------------:|:------------:|
| Monitor y dashboard | clinical-monitor-service | operation.*.data, operation.response.data |

