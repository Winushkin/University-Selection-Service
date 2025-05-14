import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import ProfileForm from './pages/ProfileForm.jsx';
import Home from './pages/Home';
import './App.css';
import EditProfileForm from "./pages/EditProfileForm.jsx";
import RegistrationForm from "./pages/RegistrationForm.jsx";
import LogInForm from "./pages/LogInForm.jsx";
import MainPage from "./pages/MainPage.jsx";

// функция маршрутизации
function App() {
    return (
        <div className="App">

            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="ProfileForm" element={<ProfileForm />} />
                <Route path="EditProfileForm" element={<EditProfileForm />} />
                <Route path="RegistrationForm" element={<RegistrationForm />} />
                <Route path="LogInForm" element={<LogInForm />} />
                <Route path="MainPage" element={<MainPage />} />
            </Routes>

        </div>
    );
}

export default App;
