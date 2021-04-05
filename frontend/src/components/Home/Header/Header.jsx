import React, { useState } from 'react';
import { Container, Row } from 'react-bootstrap';
import {
  faChalkboardTeacher, faUsers,
} from '@fortawesome/free-solid-svg-icons';

import ControlMenu from './ControlMenu/ControlMenu';
import HelpModal from './HelpModal/HelpModal';

import logo from '../../../assets/homeHeader.svg';
import './header.sass';

const Header = () => {
  const windowEnum = {
    NONE: 0, TEAM: 1, BOARDS: 2, MODAL: 3,
  };

  const [activeWindow, setActiveWindow] = useState(windowEnum.NONE);

  const handleActivate = (window) => () => {
    switch (window) {
      case activeWindow: setActiveWindow(windowEnum.NONE); break;
      case windowEnum.TEAM: setActiveWindow(windowEnum.TEAM); break;
      case windowEnum.BOARDS: setActiveWindow(windowEnum.BOARDS); break;
      case windowEnum.MODAL: setActiveWindow(windowEnum.MODAL); break;
      default: setActiveWindow(windowEnum.NONE); break;
    }
  };

  return (
    <div id="Header">
      <div className="Logo">
        <img alt="logo" src={logo} />
      </div>
      <div className="ControlsWrapper">
        <Container>
          <Row className="ControlsRow">
            <ControlMenu
              name="team"
              isActive={activeWindow === windowEnum.TEAM}
              activate={handleActivate(windowEnum.TEAM)}
              icon={faUsers}
            />

            <ControlMenu
              name="boards"
              isActive={activeWindow === windowEnum.BOARDS}
              activate={handleActivate(windowEnum.BOARDS)}
              icon={faChalkboardTeacher}
            />

            <HelpModal
              isActive={activeWindow === windowEnum.MODAL}
              activate={handleActivate(windowEnum.MODAL)}
            />
          </Row>
        </Container>
      </div>
    </div>
  );
};

export default Header;
