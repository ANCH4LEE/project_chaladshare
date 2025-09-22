import './App.css';
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import Login from './pages/login';
import Register from './pages/register';
import Home from './pages/Home';
import PostDetail from './pages/PostDetail';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<Home />} />
        <Route path="/post/:id" element={<PostDetail />} />
  
      </Routes>
    </Router>
  );
}

export default App;
