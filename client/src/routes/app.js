import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import SignIn from './signin';
import SignUp from './signup';
import Groups from './groups';
import Error from './error';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/signup" element={<SignUp />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/groups" element={<Groups />} />
        <Route index element={<Groups />} />
        <Route path="/*" element={<Error />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
