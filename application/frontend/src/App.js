import React from "react";
import styled from "styled-components";
import { createGlobalStyle } from "styled-components";
import { ClientSideOnlyMap } from "./ClientSideOnlyMap.js";
import { createTheme, ThemeProvider } from "@mui/material/styles";

const GlobalStyle = createGlobalStyle`
  body {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  code {
    font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New',
    monospace;
  }
`;

const Div = styled.div`
  text-align: center;
`;

const theme = createTheme({
  palette: {
    primary: {
      main: "#3f51b5",
    },
    secondary: {
      main: "#f44336",
    },
  },
});

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <Div className="App">
          <p>Bonjour</p>
          <ClientSideOnlyMap />
        </Div>
      </ThemeProvider>
    </>
  );
}

export default App;
