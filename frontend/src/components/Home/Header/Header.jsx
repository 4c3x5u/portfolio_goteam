import React, { useState } from 'react';
import { Container, Row } from 'react-bootstrap';
import {
  faChalkboardTeacher, faUsers,
} from '@fortawesome/free-solid-svg-icons';

import ControlMenu from './ControlMenu/ControlMenu';
import HelpModal from './HelpModal/HelpModal';

import './header.sass';
import logo from '../../../assets/homeHeader.svg';

const Header = () => {
  const [teamControlsOn, setTeamControlsOn] = useState(false);
  const [boardControlsOn, setBoardControlsOn] = useState(false);

  const toggleTeamControls = () => {
    if (!teamControlsOn) { setBoardControlsOn(teamControlsOn); }
    setTeamControlsOn(!teamControlsOn);
  };

  const toggleBoardControls = () => {
    if (!boardControlsOn) { setTeamControlsOn(boardControlsOn); }
    setBoardControlsOn(!boardControlsOn);
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
              toggle={toggleTeamControls}
              isToggled={teamControlsOn}
              icon={faUsers}
            />

            <ControlMenu
              name="boards"
              toggle={toggleBoardControls}
              isToggled={boardControlsOn}
              icon={faChalkboardTeacher}
            />

            <HelpModal />
          </Row>
        </Container>
      </div>
    </div>
  );
};

export default Header;
