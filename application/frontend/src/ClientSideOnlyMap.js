import React, { useState, useEffect } from "react";
import { MyMap } from "./Map";
export function ClientSideOnlyMap() {
  const [isClient, setIsClient] = useState(false);
  const [isSignedIn, setIsSignedIn] = useState(false); // Local signed-in state.
  useEffect(() => {
    setIsClient(true);
  }, []);
  if (!isClient) {
    return null;
  }
  return <MyMap isSignedIn={isSignedIn} setIsSignedIn={setIsSignedIn} />;
}
