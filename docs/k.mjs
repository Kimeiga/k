
// import { initializeApp } from "./firebase-app.js";
// import { getFirestore, doc, getDoc, setDoc, deleteDoc } from "./firebase-firestore.js";
import { initializeApp } from "https://www.gstatic.com/firebasejs/9.4.0/firebase-app.js";
import { getFirestore, doc, getDoc, setDoc, deleteDoc } from "https://www.gstatic.com/firebasejs/9.4.0/firebase-firestore.js";

import { getAuth, signInWithPopup, GoogleAuthProvider, signOut } from "https://www.gstatic.com/firebasejs/9.4.0/firebase-auth.js";
// Your web app's Firebase configuration
var firebaseConfig = {
	apiKey: "AIzaSyDzJSsM3BXzjrLc_oTEHOjhJGK-vpC5_F4",
	authDomain: "kiokun2.firebaseapp.com",
	projectId: "kiokun2",
	storageBucket: "kiokun2.appspot.com",
	messagingSenderId: "340492543052",
	appId: "1:340492543052:web:95ad69592608fd75ae00e9",
	measurementId: "G-RP9MNVRHNJ"
};
// Initialize Firebase
const app = initializeApp(firebaseConfig);

const db = getFirestore(app);
const auth = getAuth(app);

class GoogleSignIn extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: 'open' });
		this.render();
		auth.onAuthStateChanged((user) => this.updateUI(user));
	}

	render() {
		this.shadowRoot.innerHTML = `
            <button id="sign-in-button">Sign in with Google</button>
            <div id="user-info" style="display:none;">
                <span id="user-name"></span>
                <button id="sign-out-button">Sign out</button>
            </div>
        `;

		this.shadowRoot.querySelector("#sign-in-button").addEventListener("click", this.signIn.bind(this));
		this.shadowRoot.querySelector("#sign-out-button").addEventListener("click", this.signOut.bind(this));
	}

	signIn() {
		const provider = new GoogleAuthProvider();
		signInWithPopup(auth, provider);
	}

	signOut() {
		signOut(auth);
	}

	updateUI(user) {
		if (user) {
			this.shadowRoot.querySelector("#user-info").style.display = 'block';
			this.shadowRoot.querySelector("#user-name").textContent = `Hello, ${user.displayName}`;
			this.shadowRoot.querySelector("#sign-in-button").style.display = 'none';
		} else {
			this.shadowRoot.querySelector("#user-info").style.display = 'none';
			this.shadowRoot.querySelector("#sign-in-button").style.display = 'block';
		}
		document.dispatchEvent(new CustomEvent('auth-change', { detail: { user } }));
	}
}

window.customElements.define('google-sign-in', GoogleSignIn);


class CharacterEntry extends HTMLElement {
	constructor() {
		super();
		this.noteText = '';
		this.timeoutId = null;
		this.changeStashed = false;
		this.user = null;
		document.addEventListener('auth-change', (e) => {
			this.user = e.detail.user;
			this.render();
			if (this.user) {
				this.loadNote();
			}
		});
	}

	connectedCallback() {
		this.render();
	}

	async loadNote() {
		const entry = this.getAttribute('word');
		const docRef = doc(db, `words/${entry}`);
		const docSnap = await getDoc(docRef);
		if (docSnap.exists()) {
			this.noteText = docSnap.data().note || '';
			this.querySelector('textarea').value = this.noteText;
		}
	}

	handleInput(event) {
		this.noteText = event.target.value;
		clearTimeout(this.timeoutId);
		this.timeoutId = setTimeout(() => this.saveNote(), 1000);
		this.changeStashed = true;
	}

	handleBlur() {
		if (this.changeStashed) {
			clearTimeout(this.timeoutId);
			this.saveNote();
			this.changeStashed = false;
		}
	}

	async saveNote() {
		const entry = this.getAttribute('word');
		const docRef = doc(db, `words/${entry}`);
		if (this.noteText) {
			await setDoc(docRef, { note: this.noteText });
		} else {
			await deleteDoc(docRef);
		}
	}

	render() {
		if (this.user) {
			// Display editable textarea
			this.innerHTML = `
                    <textarea class="note" placeholder="Note...">${this.noteText}</textarea>
            `;
			const textarea = this.querySelector('textarea');
			textarea.addEventListener('input', this.handleInput.bind(this));
			textarea.addEventListener('blur', this.handleBlur.bind(this));
		} else {
			// Display paragraph
			this.innerHTML = `
                    <p>Sign in to edit</p>
            `;
		}
	}
}

window.customElements.define('word-note', CharacterEntry);
