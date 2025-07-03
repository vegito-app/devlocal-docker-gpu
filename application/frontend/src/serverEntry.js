import 'fast-text-encoding';
import React from 'react';
import ReactDOM from 'react-dom';
import ReactDOMServer from 'react-dom/server';
import App from './App';
import { StaticRouter } from 'react-router-dom/server';
import { BrowserRouter } from 'react-router-dom';
import { ServerStyleSheet } from 'styled-components'

const hydrateApp = () =>
  ReactDOM(
    <BrowserRouter>
      <App />
    </BrowserRouter>,
    document.getElementById('root')
);

export default hydrateApp;

function renderApp(url, context = {}) {  
  const sheet = new ServerStyleSheet()
  try {
    const appHtml = ReactDOMServer.renderToString(
      sheet.collectStyles(
        <StaticRouter location={url} context={context}>
          <App />
        </StaticRouter>
      )
    );

    // Extract styles before sealing the sheet
    const styles = sheet.getStyleTags()
    
    // Seal the sheet
    sheet.seal();

    return { html: appHtml, styles };
  } catch (error) {
    console.error(error);
    // Make sure to seal the sheet in case of error
    sheet.seal();
  }
}

global.renderApp = renderApp;