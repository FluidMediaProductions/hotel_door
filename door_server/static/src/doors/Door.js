import React from 'react';
import PropTypes from "prop-types";

const Door = ({id, piId, mac, number}) => (
    <tr>
        <th scope="row">{id}</th>
        <td>{piId}</td>
        <td>{mac}</td>
        <td>{number}</td>
    </tr>
);

Door.propTypes = {
    id: PropTypes.number.isRequired,
    piId: PropTypes.number,
    mac: PropTypes.string,
    number: PropTypes.number.isRequired
};

export default Door;