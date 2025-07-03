import { useState, useEffect } from "react";

export function useConfigFromBackend(configEndpointPath) {
  const [config, setConfig] = useState(null);
  useEffect(() => {
    getConfig(configEndpointPath)
      .then((result) => {
        setConfig(result);
      })
      .catch((error) => {
        console.error(
          `get config from backend '${window.location.host}${configEndpointPath}':`,
          error
        );
      });
  }, [configEndpointPath]);
  return [config];
}
// Fonction pour obtenir la configuration Firebase du serveur backend
async function getConfig(configEndpointPath) {
  const response = await fetch(configEndpointPath);
  const data = await response.json();
  return data;
}
