import React from "react";
// Import FirebaseAuth and firebase.
import { useState, useEffect } from "react";
import StyledFirebaseAuth from "react-firebaseui/StyledFirebaseAuth";
import { Modal } from "@mui/material";
import { EmailAuthProvider, GoogleAuthProvider } from "firebase/auth";
// Utilisez maintenant `getFirebaseConfig` comme vou

const uiConfig = {
  signInFlow: "popup",
  popupMode: false,
  signInOptions: [
    {
      provider: EmailAuthProvider.PROVIDER_ID,
      signInMethod: EmailAuthProvider.EMAIL_LINK_SIGN_IN_METHOD,
    },
    {
      provider: GoogleAuthProvider.PROVIDER_ID,
      // signInMethod: GoogleAuthProvider.EMAIL_LINK_SIGN_IN_METHOD,
    },
  ],
  callbacks: {
    signInSuccessWithAuthResult: function (authResult, redirectUrl) {
      console.log("callback signInSuccessWithAuthResult", authResult.user);
      console.log(
        "callback signInSuccessWithAuthResult redirect url",
        redirectUrl
      );
      // Les informations utilisateur peuvent être consultées avec authResult.user ou authResult.credential.
      // Retourner true pour continuer la redirection vers signInSuccessUrl (si configuré)
      return true;
    },
  },
};

export const SignInScreen = ({ auth, currentUser }) => {
  const [isShowingSignInModal, setIsShowingSignInModal] = useState(false);
  // Listen to the Firebase Auth state and set the local state.
  useEffect(() => {
    if (!currentUser) {
      console.log("User is logged out");
      setIsShowingSignInModal(true);
      return;
    }
    setIsShowingSignInModal(false);
  }, [currentUser]);

  return (
    <Modal
      open={isShowingSignInModal}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
    >
      <div>
        <h1>Vegito</h1>
        <p>Version {process.env.REACT_APP_VERSION}</p>
        <p>Please sign-in:</p>
        <StyledFirebaseAuth
          uiCallback={(ui) => ui.disableAutoSignIn()}
          uiConfig={uiConfig}
          firebaseAuth={auth}
        />
      </div>
    </Modal>
  );
};
export default SignInScreen;
