import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import MarkdownDisplay from "./MarkdownDisplay"; // Ajustez le chemin en fonction de l'emplacement de votre fichier

export default function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Accueil</Link>
            </li>
            <li>
              <Link to="/about">À propos</Link>
            </li>
            <li>
              <Link to="/users">Utilisateurs</Link>
            </li>
          </ul>
        </nav>
        <Switch>
          <Route path="/about">
            <About />
          </Route>
          <Route path="/users">
            <Users />
          </Route>
          <Route path="/">
            <Home />
          </Route>
          <Route path="/blog">
            <MarkdownDisplay />{" "}
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

function Home() {
  return <h2>Accueil</h2>;
}

function About() {
  return <h2>À propos</h2>;
}

function Users() {
  return <h2>Utilisateurs</h2>;
}
