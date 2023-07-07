import { useState } from "react";

export default function TestArea() {
  const [sent, setSent] = useState<string>("Nothing sent yet...");
  const [received, setReceived] = useState<string>("Nothing received yet...");
  const [outputs, setOutputs] = useState<string[][]>([]);

  function TestBroker() {
    const body = {
      method: "POST",
    };

    fetch("http://localhost:8080", body)
      .then((response) => response.json())
      .then((data) => {
        setSent("Empty post request");
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([...outputs, ["Error from Broker service", data.message]]);
        } else {
          setOutputs([
            ...outputs,
            ["Response from Broker service", data.message],
          ]);
        }
      })
      .catch((err) => {
        console.log(err);
        setOutputs([
          ...outputs,
          ["Error fetching from Broker service", err.message],
        ]);
      });
  }

  function TestAuth() {
    const payload = {
      action: "auth",
      auth: {
        email: "admin@example.com",
        password: "verysecret",
      },
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch("http://localhost:8080/handle", body)
      .then((response) => response.json())
      .then((data) => {
        setSent(JSON.stringify(payload, undefined, 4));
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([...outputs, ["Error from Auth service", data.message]]);
        } else {
          setOutputs([
            ...outputs,
            ["Response from Auth service", data.message],
          ]);
        }
      })
      .catch((err) => {
        setOutputs([
          ...outputs,
          ["Error fetching from Auth service", err.message],
        ]);
      });
  }

  function TestRabbitMQLogger() {
    const payload = {
      action: "log",
      log: {
        name: "RabbitMQ event",
        data: "Some kind of RabbitMQ data",
      },
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch("http://localhost:8080/handle", body)
      .then((response) => response.json())
      .then((data) => {
        setSent(JSON.stringify(payload, undefined, 4));
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([
            ...outputs,
            ["Error from RPC Logger service", data.message],
          ]);
        } else {
          setOutputs([
            ...outputs,
            ["Response from RPC Logger service", data.message],
          ]);
        }
      })
      .catch((err) => {
        setOutputs([
          ...outputs,
          ["Error fetching from RPC Logger service", err.message],
        ]);
      });
  }

  function TestGRPCLogger() {
    const payload = {
      action: "log",
      log: {
        name: "gRPC Event",
        data: "Some kind of gRPC data",
      },
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch("http://localhost:8080/log-grpc", body)
      .then((response) => response.json())
      .then((data) => {
        setSent(JSON.stringify(payload, undefined, 4));
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([
            ...outputs,
            ["Error from gRPC Logger service", data.message],
          ]);
        } else {
          setOutputs([
            ...outputs,
            ["Response from gRPC Logger service", data.message],
          ]);
        }
      })
      .catch((err) => {
        setOutputs([
          ...outputs,
          ["Error fetching from gRPC Logger service", err.message],
        ]);
      });
  }

  function TestMailer() {
    const payload = {
      action: "mail",
      mail: {
        from: "me@example.com",
        to: "you@example.com",
        subject: "Test Email Subject",
        message: "Hello, world!",
      },
    };

    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    const body = {
      method: "POST",
      body: JSON.stringify(payload),
      headers: headers,
    };

    fetch("http://localhost:8080/handle", body)
      .then((response) => response.json())
      .then((data) => {
        setSent(JSON.stringify(payload, undefined, 4));
        setReceived(JSON.stringify(data, undefined, 4));
        if (data.error) {
          setOutputs([...outputs, ["Error from Mailer service", data.message]]);
        } else {
          setOutputs([
            ...outputs,
            ["Response from Mailer service", data.message],
          ]);
        }
      })
      .catch((err) => {
        setOutputs([
          ...outputs,
          ["Error fetching from Mailer service", err.message],
        ]);
      });
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
            Test Auth
          </a>
          <a
            id="logBtn"
            className="btn btn-outline-secondary"
            onClick={TestRabbitMQLogger}
          >
            Test RabbitMQ Logger
          </a>
          <a
            id="logGBtn"
            className="btn btn-outline-secondary"
            onClick={TestGRPCLogger}
          >
            Test gRPC Logger
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
            <span className="text-muted">Output shows here...</span>
            {outputs.map((o) => {
              return (
                <>
                  <br></br>
                  <strong>{o[0]}</strong>: {o[1]}
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
              <span className="text-muted">{sent}</span>
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
              <span className="text-muted">{received}</span>
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}
