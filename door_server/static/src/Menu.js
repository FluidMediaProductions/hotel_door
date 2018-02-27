import React from 'react';
import PropTypes from 'prop-types';
import MenuItem from './MenuItem';

const Menu = ({pages}) => (
    <nav className="navbar navbar-expand-lg navbar-light bg-light mb-4">
        <a className="navbar-brand" href="/">Hotel door system</a>
        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#mainNavbar">
            <span className="navbar-toggler-icon" />
        </button>

        <div className="collapse navbar-collapse" id="mainNavbar">
            <ul className="navbar-nav mr-auto">
                {pages.map(page => (
                    <MenuItem key={page.id} text={page.title} link={page.link}/>
                ))}
            </ul>
        </div>
    </nav>
);

Menu.propTypes = {
    pages: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.number.isRequired,
            title: PropTypes.string.isRequired,
            link: PropTypes.string.isRequired
        }).isRequired
    ).isRequired
};

export default Menu