import React from 'react';
import { useNavigate } from 'react-router-dom';
import './MainPage.css';

const MainPage: React.FC = () => {
    const navigate = useNavigate();

    const handleLoginClick = () => {
        navigate('/login');
    };

    return (
        <div className="main-page">
            <header className="main-header">
                <h1>Welcome to Our Platform</h1>
                <button 
                    className="login-button"
                    onClick={handleLoginClick}
                >
                    Login
                </button>
            </header>
            <main className="main-content">
                <section className="features">
                    <h2>Our Features</h2>
                    <div className="feature-cards">
                        <div className="feature-card">
                            <h3>For Clients</h3>
                            <p>Find the best tutors for your needs</p>
                        </div>
                        <div className="feature-card">
                            <h3>For Tutors</h3>
                            <p>Share your knowledge and earn</p>
                        </div>
                        <div className="feature-card">
                            <h3>For Moderators</h3>
                            <p>Help maintain platform quality</p>
                        </div>
                    </div>
                </section>
            </main>
        </div>
    );
};

export default MainPage; 