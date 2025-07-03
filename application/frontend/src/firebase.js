import { useState, useEffect } from "react";
import { useConfigFromBackend } from "./backend.js";
import { initializeApp } from "firebase/app";
import { getAuth, connectAuthEmulator } from "firebase/auth";
import { getFirestore, connectFirestoreEmulator } from "firebase/firestore";
import { getFunctions, connectFunctionsEmulator } from "firebase/functions";

const configEndpoint = "/ui/config/firebase";

export function useFirebase() {
  const [auth, setAuth] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const [firestore, setFirestore] = useState(null);
  const [functions, setFunctions] = useState(null);
  const [configFromBackend] = useConfigFromBackend(configEndpoint);

  useEffect(() => {
    if (configFromBackend === null) {
      return;
    }
    const appInstance = initializeApp(configFromBackend);
    const authInstance = getAuth(appInstance);
    const firestore = getFirestore(appInstance);
    const functions = getFunctions(appInstance);

    if (window.location.hostname === "localhost") {
      connectAuthEmulator(authInstance, "http://localhost:9099");
      connectFirestoreEmulator(firestore, "http://localhost:8090");
      connectFunctionsEmulator(functions, "http://localhost:5001");
    }

    setAuth(authInstance);
    setFirestore(firestore);
    setFunctions(functions);
  }, [configFromBackend]);

  useEffect(() => {
    if (auth) {
      const unsubscribe = auth.onAuthStateChanged((user) => {
        setCurrentUser(user);
      });
      const currentUser = auth.currentUser;
      if (!!currentUser) {
        setCurrentUser(currentUser);
      }
      // Nettoyez la souscription lorsque le composant est démonté
      return () => unsubscribe();
    }
  }, [auth]);
  return [auth, currentUser, firestore, functions];
}
