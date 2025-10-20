import { useState } from "react";
import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import { useNavigate, NavLink } from "react-router-dom"
import './Header.css'

const Header = () => {
    const navigate = useNavigate();
    const [auth, setAuth] = useState(false);

    return (
        <Navbar className="custom-navbar" expand="lg">
            <Container>
                <Navbar.Brand className="brand-logo">
                    <span className="brand-icon">🎬</span>
                    <span className="brand-text">Loomi</span>
                </Navbar.Brand>

                <Navbar.Toggle aria-controls="main-navbar-nav" className="custom-toggler" />

                <Navbar.Collapse id="main-navbar-nav">
                    <Nav className="me-auto">
                        <Nav.Link as={NavLink} to="/" className="nav-link-custom">
                            <span className="nav-icon">🏠</span>
                            Home
                        </Nav.Link>
                        <Nav.Link as={NavLink} to="/recommended" className="nav-link-custom">
                            <span className="nav-icon">⭐</span>
                            Recommended
                        </Nav.Link>
                    </Nav>

                    <Nav className="ms-auto align-items-center nav-actions">
                        {auth ? (
                            <>
                                <span className="user-greeting">
                                    Hello, <strong className="user-name">Name</strong>
                                </span>

                                <Button className="btn-custom btn-logout">
                                    <span className="btn-icon">🚪</span>
                                    Logout
                                </Button>
                            </>
                        ) : (
                            <>
                                <Button
                                    className="btn-custom btn-login me-2"
                                    onClick={() => navigate("/login")}
                                >
                                    <span className="btn-icon">🔑</span>
                                    Login
                                </Button>

                                <Button
                                    className="btn-custom btn-register"
                                    onClick={() => navigate("/register")}
                                >
                                    <span className="btn-icon">✨</span>
                                    Register
                                </Button>
                            </>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    )
};
export default Header;