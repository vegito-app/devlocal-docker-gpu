import React from "react";
import { useEffect, useState, lazy, Suspense } from "react";
import { Drawer } from "@mui/material";
import { Marker, GoogleMap, LoadScript } from "@react-google-maps/api";
import { useConfigFromBackend } from "./backend.js";
import ErrorBoundary from "./ErrorBoundary";
import { useFirebase } from "./firebase.js";

const SignInScreen = lazy(() => import("./signIn")); // Remplacer "./SignInScreen" par le chemin réel vers votre composant

const libraries = ["places"];

const googlMapsConfigEndpoint = "/ui/config/googlemaps";

export function MyMapComponent({ isSignedIn, setIsSignedIn }) {
  const mapContainerStyle = {
    width: "1500px",
    height: "1400px",
  };
  const [configFromBackend] = useConfigFromBackend(googlMapsConfigEndpoint);
  const [currentPosition, setCurrentPosition] = useState();
  const [googleMapsApiKey, setGoogleMapsApiKey] = useState(null);
  const [auth, currentUser, firestore, functions] = useFirebase();
  const [showLogin, setShowLogin] = useState(false);

  useEffect(() => {
    if (configFromBackend === null) {
      return;
    }
    setGoogleMapsApiKey(configFromBackend.apiKey);
  }, [configFromBackend]);
  useEffect(() => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        setCurrentPosition({
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        });
      },
      (error) => {
        console.error(
          "Erreur lors de l'obtention de la position actuelle : ",
          error
        );
      }
    );
  }, []);

  // useEffect(() => {
  //   if (!currentUser) {
  //     setShowLogin(true);
  //     return;
  //   }
  // }, [currentUser]);

  const handleMapClick = () => {
    console.log("click map");
    if (!currentUser) {
      setShowLogin(true);
    }
  };
  if (googleMapsApiKey === null) {
    return <div> Loading...</div>;
  }
  return (
    <LoadScript googleMapsApiKey={googleMapsApiKey} libraries={libraries}>
      <GoogleMap
        id="example-map"
        mapContainerStyle={mapContainerStyle}
        zoom={15}
        center={currentPosition}
        onClick={handleMapClick}
      >
        <ErrorBoundary>
          <Suspense fallback={<div>Loading...</div>}>
            {showLogin && (
              <SignInScreen auth={auth} currentUser={currentUser} />
            )}
          </Suspense>
        </ErrorBoundary>
        {/* {showLogin ? <SignInScreen isShowing={showLogin} handleLoginClose={handleLoginClose} /> : null} */}
        {/* Autre code... */}
        {!!currentPosition && currentPosition.lat && currentPosition.lng && (
          <Marker visible={true} label="bonjour" position={currentPosition} />
        )}
        {/* You can put your map markers or other components here */}
      </GoogleMap>
    </LoadScript>
  );
}

export function MyMap({ isSignedIn, setIsSignedIn }) {
  return (
    <>
      <Drawer variant="permanent">
        <MyMapComponent isSignedIn={isSignedIn} setIsSignedIn={setIsSignedIn} />
      </Drawer>
    </>
  );
}
