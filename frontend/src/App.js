import './App.css';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';

import Login from './pages/Login';
import Register from './pages/Register';
import Home from './pages/Home';
import PostDetail from './pages/PostDetail';
import CreatePost from './pages/Createpost';
import Friends from './pages/Friends';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/post/:id" element={<PostDetail />} />
        <Route path="/newpost" element={<CreatePost />} />
        <Route path="/friends" element={<Friends />} />
  
      </Routes>
    </Router>
  );
}

export default App;
