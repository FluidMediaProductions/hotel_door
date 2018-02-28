import React from 'react'
import Menu from "./Menu"
import Doors from "./doors"
import Pis from "./pis"
import Actions from "./actions"
import {Route} from "react-router-dom";
import Home from "./home";

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
    }
];

export const paginationLength = 20;

const App = () => (
    <div>
        <Menu pages={pages} />
        {pages.map(page => (
            <Route key={page.id} exact={page.exact} path={page.link} component={page.component}/>
        ))}
    </div>
);

export default App;