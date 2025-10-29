import useAuth from "../../hooks/useAuth.jsx";
import {useState} from "react";
import {Link, useLocation, useNavigate} from "react-router-dom";
import axiosClient from '../../api/axiosConfig.js'
import logo from '../../assets/logo.svg';

const Login = () => {
    const { setAuth } = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const location = useLocation();
    const navigate = useNavigate();

    const from = location.state?.from?.pathname || '/';

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const response = await axiosClient.post('/login', { email, password });
            console.log(response.data);

            if (response.data.error) {
                setError(response.data.error);
                return;
            }
            setAuth(response.data);
            navigate(from, { replace: true });
        } catch (error) {
            console.error(error);
            setError('Invalid email or password');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-container">
            <div className="login-card">
                <div className="login-header">
                    <img src={logo} className="login-logo" alt="Logo" />
                    <h2>Sign in</h2>
                    <p>Welcome Back! Please login to your account</p>
                </div>

                {error && <div className="login-error">{error}</div>}

                <form onSubmit={handleSubmit} className="login-form">
                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input type="email" id="email" placeholder="Enter email" value={email}
                               onChange={(e) => setEmail(e.target.value)}
                               autoFocus
                               required />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Password</label>
                        <input type="password" id="password" placeholder="Enter your password" value={password}
                               onChange={(e) => setPassword(e.target.value)}
                               required />
                    </div>

                    <button type="submit" disabled={loading} className="login-btn">
                        {loading
                            ? (
                                <>
                                    <span className="spinner-border spinner-border-sm" role="status"
                                          aria-hidden="true"></span>
                                    Logging in...
                                </>
                            )
                            : (
                                'Login'
                            )
                        }
                    </button>
                </form>

                <div className="login-footer">
                    <span>Don't have an account?</span>
                    <Link to="/register">Register here!</Link>
                </div>
            </div>
        </div>
    );
}
export default Login;