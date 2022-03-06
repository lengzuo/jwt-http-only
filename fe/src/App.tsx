import React from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  const login = () => {
    fetch('http://localhost:8080/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({'name': "my name is super frontend"}),
      credentials: 'include',
    })
    .then((response) => response.json())
    //Then with the data from the response in JSON...
    .then((data) => {
      console.log('Login Success:', data);
    })
    //Then with the error genereted...
    .catch((error) => {
      console.error('Login Error:', error);
    });
  }

  const user = () => {
    fetch('http://localhost:8080/user', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*',
      },
      credentials: 'include',
    })
    .then((response) => response.json())
    //Then with the data from the response in JSON...
    .then((data) => {
      console.log('User Success:', data);
    })
    //Then with the error genereted...
    .catch((error) => {
      console.error('User Error:', error);
    });
  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        {/* <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a> */}
        <button className="App-link" onClick={login} >Login</button>
        <br/>
        <button className="App-link" onClick={user} >User</button>
      </header>
    </div>
  );
}

export default App;
