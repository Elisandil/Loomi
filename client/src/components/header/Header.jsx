
import { useNavigate, NavLink } from "react-router-dom"
import useAuth from "../../hooks/useAuth.jsx";


const Header = () => {
    const navigate = useNavigate();
    const { auth } = useAuth();

    return (
        <header className="header">
            <div className="header-left">
                <div className="logo">
                    <div className="logo-icon"></div>
                    LooMi
                </div>
                <nav className="nav">
                    <NavLink to="/" className="nav-link">Home</NavLink>
                    <NavLink to="/recommended" className="nav-link">Recommended</NavLink>
                </nav>
            </div>

            <div className="header-right">
                <div className="search-bar">
                    <input
                        type="text"
                        className="search-input"
                        placeholder="Search"
                    />
                </div>

                <button className="icon-btn">🔔</button>
                <button className="icon-btn">🌙</button>

                {auth ? (
                    <div className="user-profile">
                        <span>{auth.first_name}</span>
                        <img
                            src={`https://ui-avatars.com/api/?name=${auth.first_name}`}
                            alt="avatar"
                            className="user-avatar"
                        />
                    </div>
                ) : (
                    <>
                        <button className="view-all-btn" onClick={() => navigate("/login")}>
                            Login
                        </button>
                        <button className="view-all-btn" onClick={() => navigate("/register")}>
                            Register
                        </button>
                    </>
                )}
            </div>
        </header>
    )
};
export default Header;