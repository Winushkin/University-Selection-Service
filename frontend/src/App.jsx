import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import ProfileForm from './components/ProfileForm';
import Home from './pages/Home';
import './App.css';
import EditProfileForm from "./pages/EditProfileForm.jsx";
import RegistrationForm from "./pages/RegistrationForm.jsx";
import LoginForm from "./pages/LoginForm.jsx";

function App() {
    return (
        <div className="App">

            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="ProfileForm" element={<ProfileForm />} />
                <Route path="EditProfileForm" element={<EditProfileForm />} />
                <Route path="RegistrationForm" element={<RegistrationForm />} />
                <Route path="LoginForm" element={<LoginForm />} />
            </Routes>
        </div>
    );
}

export default App;