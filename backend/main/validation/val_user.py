from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response


# return (is_active, response)
def validate_is_active(is_active):
    is_empty_response = Response({
        'is_active': ErrorDetail(string='Is Active cannot be empty.',
                                 code='blank')
    }, 400)

    try:
        if not str(is_active):
            return None, is_empty_response
    except ValueError:
        return None, is_empty_response

    if not isinstance(is_active, bool):
        return None, Response({
            'is_active': ErrorDetail(string='Is Active must be a boolean.',
                                     code='invalid')
        }, 400)

    return is_active, None
