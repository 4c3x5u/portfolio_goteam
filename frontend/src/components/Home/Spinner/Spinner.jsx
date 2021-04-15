import React from 'react';
import { BeatLoader } from 'react-spinners';

import './spinner.sass';

const Spinner = () => (
  <div className="Spinner">
    <BeatLoader loading color="#3e6cb4" size={50} />
  </div>
);

export default Spinner;
