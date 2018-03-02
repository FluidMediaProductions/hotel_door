import React, {Component, Fragment} from 'react'
import Menu from "./Menu"
import Doors from "./doors"
import Pis from "./pis"
import Actions from "./actions"
import {Route} from "react-router-dom";
import Home from "./home";
import Login from "./login";
import {delJWT, JWTIsValid} from './auth';

let checkJWT = null;

const Logout = () => {
    delJWT();
    if (typeof checkJWT === "function") {
        checkJWT();
    }
    return null;
};

const pages = [
    {
        id: 0,
        title: "Home",
        link: "/",
        exact: true,
        component: Home
    },
    {
        id: 1,
        title: "Doors",
        link: "/doors",
        component: Doors,
        exact: false
    },
    {
        id: 2,
        title: "Pis",
        link: "/pis",
        component: Pis,
        exact: false
    },
    {
        id: 3,
        title: "Actions",
        link: "/actions",
        component: Actions,
        exact: false
    },
    {
        id: 100,
        title: "Logout",
        link: "/logout",
        component: Logout,
        exact: false
    }
];

export const paginationLength = 20;

class App extends Component {
    constructor(props) {
        super(props);

        this.checkJWT = this.checkJWT.bind(this);
        this.state = {
            JWTIsValid: true
        }
    }

    componentDidMount() {
        this.checkJWT();
        checkJWT = this.checkJWT;
        this.timer = setInterval(this.checkJWT, 5000);
    }

    componentWillUnmount() {
        clearInterval(this.timer);
        checkJWT = null;
    }

    checkJWT() {
        const self = this;
        JWTIsValid(function (resp) {
            self.setState({
                JWTIsValid: resp
            })
        })
    }

    render() {
        return this.state.JWTIsValid ? (
            <Fragment>
                <Menu pages={pages}/>
                {pages.map(page => (
                    <Route key={page.id} exact={page.exact} path={page.link} component={page.component}/>
                ))}
            </Fragment>
        ):(
            <Login onLogin={this.checkJWT}/>
        )
    }
}

export default App;