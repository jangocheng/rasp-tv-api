angular.module('simongeeks.filters', [])
	// custom filter to filter arrays by regex
	.filter('filterRegex', function() {
		return function(input, regex) {

			if (angular.isUndefined(regex) || !angular.isString(regex) || regex.length === 0)
				return input;

			if (!angular.isArray(input)) return input;

			// if the regex is invalid just return unfiltered input
			try {
				var regexObj = new RegExp(regex, 'i');
			} catch (e) {
				return input;
			}

			var compare = function(value) {
				if (angular.isUndefined(value) || value === null) return false;
				var i;
				switch (typeof value) {
					case 'object':
						var keys = Object.keys(value);
						for (i = 0; i < keys.length; i++) {
							if (compare(value[keys[i]])) return true;
						}
						break;
					case 'array':
						for (i = 0; i < value.length; i++) {
							if (compare(value[i])) return true;
						}
						break;
					case 'number':
					case 'string':
						var str = '' + value;
						return regexObj.test(str);
				}

				return false;
			};

			return input.filter(compare);
		};
	});