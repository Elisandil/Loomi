import {useEffect, useState} from "react";
import {Link, useNavigate} from "react-router-dom";
import axiosClient from '../../api/axiosConfig.js'
import logo from '../../assets/logo.svg';

const Register = () => {
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [favouriteGenres, setFavouriteGenres] = useState([]);
    const [genres, setGenres] = useState([])

    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate()

    const handleGenreChange = (e) => {
        const options = Array.from(e.target.selectedOptions);
        setFavouriteGenres(options.map(opt => ({
            genre_id: Number(opt.value),
            genre_name: opt.label
        })));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        const defaultRole ="USER";

        if (password !== confirmPassword) {
            setError('Passwords do not match.');
            return;
        }
        setLoading(true);
        try {
            const payload = {
                first_name: firstName,
                last_name: lastName,
                email,
                password,
                role: defaultRole,
                favourite_genres: favouriteGenres
            };
            const response = await axiosClient.post('/register', payload);
            if (response.data.error) {
                setError(response.data.error);
                return;
            }
            navigate('/login', { replace: true });
        } catch (err) {
            setError('Registration failed. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const fetchGenres = async() => {
            try {
                const response = await axiosClient.get('genres');
                setGenres(response.data);
            } catch (error) {
                console.error('Error fetching genres: ', error);
            }
        }
        fetchGenres();
    }, []);

    return (
        <div className="login-container">
            <div className="login-card" style={{ maxWidth: '600px' }}>
                <div className="login-header">
                    <img src={logo} className="login-logo" alt="Logo" />
                    <h2>Create Account</h2>
                    <p>Join Loomi and start streaming</p>
                </div>

                {error && <div className="login-error">{error}</div>}

                <form onSubmit={handleSubmit} className="login-form">
                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <div className="form-group" style={{ flex: 1 }}>
                            <label htmlFor="firstName">First Name</label>
                            <input type="text" id="firstName" placeholder="Enter first name" value={firstName}
                                   onChange={(e) => setFirstName(e.target.value)}
                                   required />
                        </div>

                        <div className="form-group" style={{ flex: 1 }}>
                            <label htmlFor="lastName">Last Name</label>
                            <input type="text" id="lastName" placeholder="Enter last name" value={lastName}
                                   onChange={(e) => setLastName(e.target.value)}
                                   required />
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input type="email" id="email" placeholder="Enter email" value={email}
                               onChange={(e) => setEmail(e.target.value)}
                               required />
                    </div>

                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <div className="form-group" style={{ flex: 1 }}>
                            <label htmlFor="password">Password</label>
                            <input type="password" id="password" placeholder="Enter password" value={password}
                                   onChange={(e) => setPassword(e.target.value)}
                                   required />
                        </div>

                        <div className="form-group" style={{ flex: 1 }}>
                            <label htmlFor="confirmPassword">Confirm Password</label>
                            <input type="password" id="confirmPassword" placeholder="Confirm password"
                                   value={confirmPassword}
                                   onChange={(e) => setConfirmPassword(e.target.value)}
                                   style={{
                                       borderColor: confirmPassword && password !== confirmPassword
                                           ? '#dc3545'
                                           : 'rgba(255, 107, 53, 0.2)'
                                   }}
                                   required />
                            {confirmPassword && password !== confirmPassword && (
                                <small style={{ color: '#ff6b6b', fontSize: '0.85rem', marginTop: '0.25rem',
                                    display: 'block' }}>
                                    Passwords do not match
                                </small>
                            )}
                        </div>
                    </div>

                    <div className="form-group">
                        <label htmlFor="genres">Favourite Genres</label>
                        <select id="genres" multiple value={favouriteGenres.map(g => String(g.genre_id))}
                                onChange={handleGenreChange}
                                style={{background: 'rgba(255, 255, 255, 0.05)',
                                    border: '1px solid rgba(255, 107, 53, 0.2)', borderRadius: '8px',
                                    padding: '0.75rem 1rem', color: 'var(--text-primary)', fontSize: '0.95rem',
                                    minHeight: '120px'
                                }}>
                            {genres.map(genre => (
                                <option key={genre.genre_id} value={genre.genre_id}>
                                    {genre.genre_name}
                                </option>
                            ))}
                        </select>
                        <small style={{ color: 'var(--text-secondary)', fontSize: '0.85rem', marginTop: '0.25rem',
                            display: 'block' }}>
                            Hold Ctrl (Windows) or Cmd (Mac) to select multiple genres
                        </small>
                    </div>

                    <button type="submit" disabled={loading} className="login-btn">
                        {loading
                            ? (
                                <>
                                    <span className="spinner-border spinner-border-sm" role="status"
                                          aria-hidden="true"></span>
                                    Registering...
                                </>
                            )
                            : (
                                'Register'
                            )
                        }
                    </button>
                </form>

                <div className="login-footer">
                    <span>Already have an account?</span>
                    <Link to="/login">Login here!</Link>
                </div>
            </div>
        </div>
    );
}
export default Register;