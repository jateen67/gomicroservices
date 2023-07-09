import { useState } from "react";

export default function TestArea() {
  const [sent, setSent] = useState<string>("Nothing sent yet...");
  const [received, setReceived] = useState<string>("Nothing received yet...");
  const [outputs, setOutputs] = useState<string[][]>([]);

  function fetchData(url: string, payload: object, serviceName: string) {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch(url, body)
      .then((response) => response.json())
      .then((data) => {
        setSent(JSON.stringify(payload, undefined, 4));
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([
            [
              `Error from ${serviceName} service`,
              data.message,
              new Date().toString(),
            ],
          ]);
        } else {
          setOutputs([
            [
              `Response from ${serviceName} service`,
              data.message,
              new Date().toString(),
            ],
          ]);
        }
      })
      .catch((err) => {
        setOutputs([
          [
            `Error fetching from ${serviceName} service`,
            err.message,
            new Date().toString(),
          ],
        ]);
      });
  }

  function TestBroker() {
    const body = {
      Content: "Empty post request",
    };

    fetchData("http://localhost:8080", body, "Broker");
  }

  function TestAuth() {
    const payload = {
      action: "auth",
      auth: {
        email: "admin@example.com",
        password: "verysecret",
      },
    };

    fetchData("http://localhost:8080/handle", payload, "Authentication");
  }

  function TestGRPCLogger() {
    const payload = {
      action: "log",
      log: {
        name: "gRPC Event",
        data: "Some kind of gRPC data",
      },
    };

    fetchData("http://localhost:8080/log-grpc", payload, "gRPC Logger");
  }

  function TestRabbitMQLogger() {
    const payload = {
      action: "log",
      log: {
        name: "RabbitMQ event",
        data: "Some kind of RabbitMQ data",
      },
    };

    fetchData("http://localhost:8080/handle", payload, "RabbitMQ Logger");
  }

  function TestMailer() {
    const payload = {
      action: "mail",
      mail: {
        from: "me@example.com",
        to: "you@example.com",
        subject: "Test Email Subject",
        message: "Hello, world! This is my email",
      },
    };

    fetchData("http://localhost:8080/handle", payload, "Mailer Service");
  }

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="mt-5">Microservices in Go</h1>
          <hr></hr>
          <a
            id="brokerBtn"
            className="btn btn-outline-secondary"
            onClick={TestBroker}
          >
            Test Broker
          </a>
          <a
            id="authBrokerBtn"
            className="btn btn-outline-secondary"
            onClick={TestAuth}
          >
            Test Authentication
          </a>
          <a
            id="logGBtn"
            className="btn btn-outline-secondary"
            onClick={TestGRPCLogger}
          >
            Test gRPC Logger
          </a>
          <a
            id="logBtn"
            className="btn btn-outline-secondary"
            onClick={TestRabbitMQLogger}
          >
            Test RabbitMQ Logger
          </a>
          <a
            id="mailBtn"
            className="btn btn-outline-secondary"
            onClick={TestMailer}
          >
            Test Mailer
          </a>

          <div
            id="output"
            className="mt-5"
            style={{ outline: "1px solid silver", padding: "2em" }}
          >
            {outputs.length === 0 ? (
              <>
                <span className="text-muted">Output shows here...</span>
              </>
            ) : (
              <></>
            )}
            {outputs.map((o) => {
              return (
                <>
                  <strong className="text-success">Started</strong>
                  <br></br>
                  <i>Sending request...</i>
                  <br></br>
                  <strong>{o[0]}</strong>: {o[1]}
                  <br></br>
                  <strong className="text-danger">Ended</strong>: {o[2]}
                </>
              );
            })}
          </div>
        </div>
      </div>
      <div className="row">
        <div className="col">
          <h4 className="mt-5">Sent</h4>
          <div
            className="mt-1"
            style={{ outline: "1px solid silver", padding: "2em" }}
          >
            <pre id="payload">
              <span style={{ fontWeight: "bold" }}>{sent}</span>
            </pre>
          </div>
        </div>
        <div className="col">
          <h4 className="mt-5">Received</h4>
          <div
            className="mt-1"
            style={{ outline: "1px solid silver", padding: "2em" }}
          >
            <pre id="received">
              <span style={{ fontWeight: "bold" }}>{received}</span>
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}
