import status

from .val_custom import CustomAPIException


def validate_column_id(column_id):
    if not column_id:
        raise CustomAPIException('column_id',
                                 'Column ID cannot be empty.',
                                 status.HTTP_400_BAD_REQUEST)

    try:
        int(column_id)
    except ValueError:
        raise CustomAPIException('column_id',
                                 'Column ID must be a number.',
                                 status.HTTP_400_BAD_REQUEST)
