# LIMS Client

## Overview
This project is a simplified implementation for sending HL7 ORM messages. It supports serialization and transmission of ORM^O01 messages via a REST API endpoint.

Currently, the application only supports the `SendORM` operation.

---

## Supported Endpoint

### Send ORM

**Endpoint:**
```
POST http://localhost:8080/v1/hl-seven/orm
```

**Description:**
This endpoint accepts a JSON request representing an HL7 ORM message and processes it into an HL7-compatible string for further handling.

---

## Example Request

```json
{
  "orm": {
    "msh": {
      "field_separator": "|",
      "sending_application": "EHRSystem",
      "sending_facility": "EHRFac",
      "receiving_application": "LabSystem",
      "receiving_facility": "LabFac",
      "message_date_time": "202312050830",
      "message_type": "ORM^O01",
      "message_control_id": "12345",
      "processing_id": "P",
      "version_id": "2.3"
    },
    "pid": {
      "patient_id": "123456",
      "patient_name": "Doe^John",
      "date_of_birth": "19700101",
      "gender": "M",
      "address": "123 Main St^^Anytown^CA^12345",
      "phone_number": "(555)555-5555"
    },
    "orc": {
      "order_control": "NW",
      "order_id": "1234",
      "placer_order_number": "5678",
      "order_status": "1^Routine",
      "order_priority": "",
      "order_date_time": "202312050830"
    },
    "obr": {
      "set_id": "1",
      "placer_order_number": "1234",
      "universal_service_id": "BMP^Basic Metabolic Panel",
      "request_date_time": "202312050830"
    }
  }
}
```

---

## Example Response

**Response (200 OK):**
```json
{
  "ack": {
    "msh": {
      "field_separator": "",
      "sending_application": "",
      "sending_facility": "",
      "receiving_application": "",
      "receiving_facility": "",
      "message_date_time": "",
      "message_type": "",
      "message_control_id": "",
      "processing_id": "",
      "version_id": ""
    },
    "msa": {
      "acknowledgment_code": "MSA",
      "message_control_id": "AA",
      "text_message": "Initial Connection Established"
    }
  }
}
```

---

## How It Works

1. **JSON Parsing:**
    - The JSON request is deserialized into Go structs (`SendORMRequest` and nested `ORM`, `MSH`, `PID`, `ORC`, `OBR`).

2. **HL7 Serialization:**
    - The `Serialize` function processes the structs and converts them into an HL7-compatible string.

3. **Response:**
    - The serialized HL7 string is returned in the response for integration with downstream systems.

---

## Running the Service

1. **Start the Server:**
    - Run the application to start the HTTP server on `localhost:8080`.

2. **Test the Endpoint:**
    - Use tools like `Postman` or `curl` to send a `POST` request to `http://localhost:8080/v1/hl-seven/orm` with the example JSON payload.

---

## Notes

- Only ORM^O01 messages are currently supported.
- Ensure the request payload strictly follows the expected JSON schema.
- The mock server currently returns a static response and needs to be updated for dynamic processing.

---

## Example with cURL

```sh
curl -X POST http://localhost:8080/v1/hl-seven/orm \\
-H "Content-Type: application/json" \\
-d '{
  "orm": {
    "msh": {
      "field_separator": "|",
      "sending_application": "EHRSystem",
      "sending_facility": "EHRFac",
      "receiving_application": "LabSystem",
      "receiving_facility": "LabFac",
      "message_date_time": "202312050830",
      "message_type": "ORM^O01",
      "message_control_id": "12345",
      "processing_id": "P",
      "version_id": "2.3"
    },
    "pid": {
      "patient_id": "123456",
      "patient_name": "Doe^John",
      "date_of_birth": "19700101",
      "gender": "M",
      "address": "123 Main St^^Anytown^CA^12345",
      "phone_number": "(555)555-5555"
    },
    "orc": {
      "order_control": "NW",
      "order_id": "1234",
      "placer_order_number": "5678",
      "order_status": "1^Routine",
      "order_priority": "",
      "order_date_time": "202312050830"
    },
    "obr": {
      "set_id": "1",
      "placer_order_number": "1234",
      "universal_service_id": "BMP^Basic Metabolic Panel",
      "request_date_time": "202312050830"
    }
  }
}'
```