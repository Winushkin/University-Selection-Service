import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import ProfileForm from './components/ProfileForm';
import Home from './pages/Home';
import './App.css';
import EditProfileForm from "./pages/EditProfileForm.jsx";

function App() {
    return (
        <div className="App">

            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="ProfileForm" element={<ProfileForm />} />
                <Route path="EditProfileForm" element={<EditProfileForm />} />
            </Routes>
        </div>
    );
}

export default App;