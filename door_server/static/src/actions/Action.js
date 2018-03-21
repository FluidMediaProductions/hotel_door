import React from 'react';
import PropTypes from "prop-types";

const Action = ({id, type, piMac, complete, success}) => {
    let completeText, successText = null;
    if (complete) {
        completeText = <span className="text-success">Complete</span>;
        if (success) {
            successText = <span className="text-success">Successful</span>;
        } else {
            successText = <span className="text-danger">Failed</span>;
        }
    } else {
        completeText = <span className="text-warning">Waiting</span>;
    }
    return (
        <tr>
            <th scope="row">{id}</th>
            <td>{type}</td>
            <td>{piMac}</td>
            <td>{completeText}</td>
            <td>{successText}</td>
        </tr>
    );
};

Action.propTypes = {
    id: PropTypes.number.isRequired,
    type: PropTypes.string.isRequired,
    piMac: PropTypes.string.isRequired,
    complete: PropTypes.bool.isRequired,
    success: PropTypes.bool.isRequired,
};

export default Action;