import './App.css';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import Login from './pages/login';
import Register from './pages/register';
import Home from './pages/Home';
import PostDetail from './pages/PostDetail';
import CreatePost from './pages/Createpost';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/post/:id" element={<PostDetail />} />
        <Route path="/newpost" element={<CreatePost />} />
  
      </Routes>
    </Router>
  );
}

export default App;
