export function Carouselify(imageSrcs: string[]) {
	return imageSrcs.map((imgSrc) => {
		return {
			src: imgSrc,
			alt: 'you could see something there',
			title: 'yes'
		};
	});
}
