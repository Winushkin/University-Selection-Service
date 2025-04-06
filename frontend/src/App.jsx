import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import ProfileForm from './components/ProfileForm';
import Home from './pages/Home';
import './App.css';

function App() {
    return (
        <div className="App">

            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="ProfileForm" element={<ProfileForm />} />
            </Routes>
        </div>
    );
}

export default App;