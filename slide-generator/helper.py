def validateSectionInput(sections):
    for i, section in enumerate(sections):
        if not isinstance(section, dict):
            raise ValueError(f"Section {i} is not an object")
        if "slides" not in section or not isinstance(section["slides"], list):
            raise ValueError(f"Section {i} missing 'slides' or 'slides' is not an array")
        for j, slide in enumerate(section["slides"]):
            if not isinstance(slide, dict):
                raise ValueError(f"Slide {j} in section {i} is not an object")
            if "title" not in slide or not isinstance(slide["title"], str):
                raise ValueError(f"Slide {j} in section {i} missing 'title' or 'title' is not a string")
            if "key_points" not in slide or not isinstance(slide["key_points"], list):
                raise ValueError(f"Slide {j} in section {i} missing 'key_points' or 'key_points' is not an array")
            if "layout" not in slide or not isinstance(slide["layout"], str):
                raise ValueError(f"Slide {j} in section {i} missing 'layout' or 'layout' is not a string")
            