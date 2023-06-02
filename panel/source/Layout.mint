component Layout {
  style layout {
    display: flex;
    flex-direction: row;
    height: 100vh;
    width: 100vw;
  }

  style layout-nav {
    margin: 2em;
    list-style-type: none;
    border-right: 1px solid;
    padding: 1em;
    width: 24%;
  }

  style layout-pane {
    margin: 1em;
    width: 86%;
  }

  style layout-nav-item {
    list-style-type: none;
    padding: 1em;

    &:hover {
      background-color: #1a1a1a;
      color: #e0e0e0;
    }
  }

  style layout-nav-logo {
    list-style-type: none;
    padding: 1em;
    align-content: center;
    display: flex;
    flex-direction: column;
  }

  property children : Array(Html) = []

  fun render : Html {
    <section::layout>
      <section::layout-nav>
        <ul>
          <li::layout-nav-logo>
            <img
              src={@asset(./favicon.png)}
              style="margin:auto;"/>

            <h3 style="margin: 10px 0px; text-align:center;">
              "Providence"
            </h3>
          </li>

          <li::layout-nav-item>"Dashboard"</li>
          <li::layout-nav-item>"Clients"</li>
          <li::layout-nav-item>"Tasks"</li>
        </ul>
      </section>

      <section::layout-pane>
        <{ children }>
      </section>
    </section>
  }
}
