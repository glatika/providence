record ClientInfo {
  hwid : String,
  operatingSystem : String using "operating_system",
  signature : String,
  variant : String
}

enum Status {
  Initial
  Loading
  Error(String)
  Ok(Array(ClientInfo))
}

store ClientOnline {
  state status : Status = Status::Initial

  fun load : Promise(Void) {
    next { status: Status::Loading }

    let req =
      "http://localhost:4877/rpc/client-page-1-10"
      |> Http.get()
      |> Http.send()

    case await req {
      Result::Err => next { status: Status::Error("Something went wrong with the request.") }

      Result::Ok(response) =>
        {
          let result =
            response.bodyString
            |> Json.parse()

          case result {
            Result::Err => next { status: Status::Error("The data is not what is expected.") }

            Result::Ok(object) =>
              case decode object as Array(ClientInfo) {
                Result::Ok(r) => next { status: Status::Ok(r) }
                Result::Err(msg) => next { status: Status::Error("Invalid JSON data." + Object.Error.toString(msg)) }
              }
          }
        }
    }
  }
}

routes {
  * {
    ClientOnline.load()
  }
}

component Dashboard {
  connect ClientOnline exposing { status }

  style dashboard-let {
    flex-direction: row;
    display: flex;
  }

  style dashboard-item {
    margin: 2em;
    padding: 1em;
    border-bottom: 1px solid black;
  }

  style active-client-table {
    border-collapse: collapse;
    text-align: center;

    width: 80%;
    margin-right: auto;
  }

  style active-client-row {
    padding: 1em;
    border-top: 1px solid;
    max-width: 6em;
  }

  style active-client-row-hover {
    &:hover {
      background-color: #1a1a1a;
      color: #e0e0e0;
    }
  }

  style active-client-row-head {
    padding: 1em;
  }

  fun render : Html {
    <section>
      <section::dashboard-let id="top-dashboard">
        <section::dashboard-item>
          case status {
            Status::Ok(clients) =>
              <section style="display: flex;flex-direction: row;">
                <h1 style="font-size: 3em;margin: 5px;">
                  <{
                    if Array.size(clients) < 10 {
                      "0"
                    } else {
                      ""
                    }
                  }>

                  <{
                    Array.size(clients)
                    |> Number.toString()
                  }>
                </h1>

                <h3>
                  <{
                    if Array.size(clients) < 10 {
                      "/0"
                    } else {
                      "/"
                    }
                  }>

                  <{
                    Array.size(clients)
                    |> Number.toString()
                  }>
                </h3>
              </section>

            =>
              <section style="display: flex;flex-direction: row;">
                <h1 style="font-size: 3em;margin: 5px;">
                  "--"
                </h1>

                <h3>"/--"</h3>
              </section>
          }

          <p>"Online client"</p>
        </section>

        case status {
          Status::Ok(clients) =>
            <section::dashboard-item>
              <h1 style="font-size: 3em;margin: 5px;">
                <{
                  Array.size(clients)
                  |> Number.toString()
                }>
              </h1>

              <p>"Total client"</p>
            </section>

          =>
            <section::dashboard-item>
              <h1 style="font-size: 3em;margin: 5px;">
                "-"
              </h1>

              <p>"Total client"</p>
            </section>
        }
      </section>

      <section::dashboard-let id="chart-dashboard">
        case status {
          Status::Initial => <section/>
          Status::Loading => <p>"Loading the data..."</p>

          Status::Error(msg) =>
            <p>
              <{ msg }>
            </p>

          Status::Ok(clients) =>
            <table::active-client-table>
              <caption>"Currently Active Client"</caption>

              <thead>
                <tr>
                  <th::active-client-row-head>"Hwid"</th>
                  <th::active-client-row-head>"Operating System"</th>
                  <th::active-client-row-head>"Variant"</th>
                  <th::active-client-row-head>"Signature"</th>
                </tr>
              </thead>

              <tbody>
                for client of clients {
                  <tr::active-client-row-hover>
                    <td::active-client-row>
                      <{ client.hwid }>
                    </td>

                    <td::active-client-row>
                      <{ client.operatingSystem }>
                    </td>

                    <td::active-client-row>
                      <{ client.variant }>
                    </td>

                    <td::active-client-row>
                      <{ client.signature }>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
        }
      </section>
    </section>
  }
}
