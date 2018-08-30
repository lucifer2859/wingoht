import igt_base_template


class APNTemplate(igt_base_template.BaseTemplate):
    def __init__(self):
        igt_base_template.BaseTemplate.__init__(self)

    def getTemplateId(self):
        """templateid support,you do not need to call this function explicitly"""
        return 5
