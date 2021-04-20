import React, { useContext } from 'react';
import PropTypes from 'prop-types';
import {
  Container, Row, OverlayTrigger, Tooltip,
} from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUsers, faSignOutAlt, faUser } from '@fortawesome/free-solid-svg-icons';

import BoardsControls from './Controls/Boards/BoardsControls';
import TeamControls from './Controls/Team/TeamControls';
import HelpToggler from './HelpToggler/HelpToggler';
import window from '../../../misc/window';
import AppContext from '../../../AppContext';

import logo from '../home.svg';
import './header.sass';

const Header = ({ activeWindow, handleActivate }) => {
  const { user, loadBoard } = useContext(AppContext);

  const logout = () => {
    sessionStorage.removeItem('username');
    sessionStorage.removeItem('auth-token');
    loadBoard();
  };

  return (
    <div id="Header">
      <div className="Logo">
        <img alt="logo" src={logo} />
      </div>
      <div className="ControlsWrapper">
        <Container>
          <Row className="ControlsRow">
            <button className="LogoutButton" onClick={logout} type="button">
              <FontAwesomeIcon
                className="Icon fa-flip-horizontal"
                icon={faSignOutAlt}
              />
            </button>

            <TeamControls
              isActive={activeWindow === window.TEAM}
              handleActivate={handleActivate(window.TEAM)}
              handleCreate={handleActivate(window.INVITE_MEMBER)}
              handleDelete={handleActivate(window.DELETE_MEMBER)}
              icon={faUsers}
            />

            <BoardsControls
              isActive={activeWindow === window.BOARDS}
              handleActivate={handleActivate(window.BOARDS)}
              handleCreate={handleActivate(window.CREATE_BOARD)}
              handleDelete={handleActivate(window.DELETE_BOARD)}
              handleEdit={handleActivate(window.EDIT_BOARD)}
            />

            <HelpToggler toggle={handleActivate(window.HELP)} />

            <OverlayTrigger
              placement="bottom"
              overlay={(
                <Tooltip id="UserTooltip">
                  {`Logged in as ${user.username}.`}
                </Tooltip>
              )}
            >
              <button className="UserButton" type="button">
                <FontAwesomeIcon
                  className="Icon fa-flip-horizontal"
                  icon={faUser}
                />
              </button>
            </OverlayTrigger>
          </Row>
        </Container>
      </div>
    </div>
  );
};

Header.propTypes = {
  activeWindow: PropTypes.number.isRequired,
  handleActivate: PropTypes.func.isRequired,
};

export default Header;
