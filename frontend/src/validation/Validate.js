import { capFirstLetterOf } from '../misc/util';

const Validate = {
  requiredString: (field, maxLength, minLength) => (value) => {
    const fieldName = capFirstLetterOf(field);
    if (!value) {
      return `${fieldName} cannot be empty.`;
    }
    if (maxLength && value.length > maxLength) {
      return `${fieldName} cannot be longer than ${maxLength} characters.`;
    }
    if (minLength && value.length < minLength) {
      return `${fieldName} cannot be longer than ${minLength} characters.`;
    }
    return '';
  },
};

export default Validate;
