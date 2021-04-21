/* eslint-disable
jsx-a11y/no-static-element-interactions,
jsx-a11y/click-events-have-key-events */

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from 'react-bootstrap';

import './help.sass';
import logo from './help.svg';

const Help = ({ toggleOff }) => (
  <div className="Help" onClick={toggleOff}>
    <div className="Body" onClick={(e) => e.stopPropagation()}>
      <div className="HeaderWrapper">
        <img className="Header" alt="logo" src={logo} />
      </div>

      <h1>Summary</h1>
      <p>
        In both the control menus and the columns, you can click the
        plus icon to add a new item.
        <br />
        Also, right click on any item to view additional controls for it.
        Finally, you can toggle this modal off
        <br />
        by clicking anywhere outside its body.
      </p>

      <h1>Team Controls</h1>
      <ol>
        <li>
          Click the TEAM button located on the control bar to toggle the
          Team controls menu.
        </li>
        <li>
          The team members that are included in the currently displayed
          board and therefore can see
          <br />
          and interact with its items will have a yellow tick on both
          sides of their name.
        </li>
        <li>
          Click a member to include or exclude them from the current
          board.
        </li>
        <li>
          You can click the plus icon at the bottom of the menu to invite
          new members to your team.
        </li>
      </ol>

      <h1>Boards Controls</h1>
      <ol>
        <li>
          Click the BOARDS button located on the control bar to toggle
          the Boards controls menu.
        </li>
        <li>
          The currently displayed board will have arrows pointing to its
          name on both of its name.
        </li>
        <li>Click a board to load it on page.</li>
        <li>
          You can click the plus icon at the bottom of the list to create
          a new board.
        </li>
      </ol>

      <h1>Task Controls</h1>
      <ol>
        <li>
          Click the plus icon inside the INBOX column to create a new
          task.
        </li>
        <li>
          Give the task a name. You can optionally give it a description
          and/or associated subtasks.
        </li>
        <li>
          Drag and drop tasks from one column to another in order to move
          them around.
        </li>
        <li>
          Right click tasks to edit or delete it, or assign a teammate
          to it.
        </li>
      </ol>
      <Button
        className="Button"
        type="button"
        aria-label="cancel"
        onClick={toggleOff}
      >
        OK
      </Button>
    </div>
  </div>
);

Help.propTypes = { toggleOff: PropTypes.func.isRequired };

export default Help;
