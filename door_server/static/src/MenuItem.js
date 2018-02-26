import React from 'react';
import PropTypes from 'prop-types';
import {Link} from "react-router-dom";

const MenuItem = ({ link, text }) => (
    <li className="nav-item">
        <Link className="nav-link" to={link}>{text}</Link>
    </li>
);

MenuItem.propTypes = {
    link: PropTypes.string.isRequired,
    text: PropTypes.string.isRequired
};

export default MenuItem
