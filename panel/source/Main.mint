component Main {
  style app {
    flex-direction: column;
    display: flex;

    background-color: #e0e0e0;
    color: #1f1f1f;
    height: 100vh;
    width: 100vw;

    font-family: Consolas, monospace;
  }

  fun render : Html {
    <div::app>
      <Layout>
        <Dashboard/>
      </Layout>
    </div>
  }
}
