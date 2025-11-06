import './App.css';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';

import Login from './pages/Login';
import Register from './pages/Register';
import Home from './pages/Home';
import PostDetail from './pages/PostDetail';
import CreatePost from './pages/Createpost';
import Friends from './pages/Friends';
import Profile from './pages/Profile';
import ForgotPassword from './pages/ForgotPassword';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/posts/:id" element={<PostDetail />} />
        <Route path="/newpost" element={<CreatePost />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/profile/:id" element={<Profile />} />
        <Route path="/forgot_password" element={<ForgotPassword />} />
  
      </Routes>
    </Router>
  );
}

export default App;
